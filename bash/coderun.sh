# 打印日志
echo "开始运行程序"

# 环境变量
export QT_QPA_PLATFORM=offscreen
export LD_LIBRARY_PATH=/opt/qt5.9.2/lib:$LD_LIBRARY_PATH

# 列出当前目录下的所有文件
ls

# 导入 GDB 脚本
cat > crash.gdb <<'EOF'
set pagination off
set confirm off
set disable-randomization off
set print thread-events off

# 捕获常见程序崩溃信号
handle SIGSEGV SIGABRT SIGBUS SIGILL SIGFPE stop print nopass

run

if $_isvoid($_exitcode)
    if $_isvoid($_exitsignal)
        # 程序仍停在信号现场，可以读取堆栈
        set $caught_signal = $_siginfo.si_signo

        echo \n========== CRASH BACKTRACE ==========\n
        thread apply all bt 100
        echo ========== END BACKTRACE ==========\n

        # 立即终止被 GDB 暂停的程序
        kill
        quit 128 + $caught_signal
    else
        # SIGKILL/OOM：程序已经死亡，无法再获取堆栈
        quit 128 + $_exitsignal
    end
else
    # 正常退出或程序主动返回非零值
    quit $_exitcode
end
EOF

# 运行程序(gdb 调试)
gdb \
--return-child-result \
--batch \
-ex "set args --offscreen --exam --freq=MAX --record --ResultLogFile={RunResultFileName} --RecordOutputFile={RecordOutputFileName}" \
-x crash.gdb \
./newAOE \
> "{workdir}/crash/crash_{id}_{indices}.log" 2>&1

# 退出码
EXIT_CODE=$?

# 打印日志
echo "程序运行结束"
echo "程序退出码: $EXIT_CODE"

# 退出码分类判断
if [ $EXIT_CODE -eq 0 ]; then
    echo "程序正常退出"
elif [ $EXIT_CODE -eq 137 ]; then
    # 128 + 9(SIGKILL),通常为 OOM 被内核杀死
    echo "检测到程序被 SIGKILL 杀死(退出码 137),疑似 OOM(内存不足)"
    echo "可执行以下命令确认是否为内核 OOM killer 所为:"
    echo "  dmesg | grep -i 'killed process'"
elif [ $EXIT_CODE -gt 128 ]; then
    # 其他信号终止,如 139=段错误(SIGSEGV)、134=SIGABRT
    echo "检测到程序被信号终止,信号号: $((EXIT_CODE - 128))"
else
    echo "检测到程序异常退出(非信号)"
fi

# 退出程序
exit $EXIT_CODE
