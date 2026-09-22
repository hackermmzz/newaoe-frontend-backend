#先copy一下newaoe目录
cp -r newaoe newaoe_copy
#修复所有文件的大小写
fixcase -f ./newaoe_copy
#切换到newaoe_copy目录下
cd newaoe_copy
#添加qt5.9.2的bin目录到PATH
export PATH="$PATH:/opt/qt5.9.2/bin/"
###############################################
# 编译Release版本
##############################################
if ! qmake \
    CONFIG+=release \
    CONFIG-=debug \
    "QMAKE_CXXFLAGS_RELEASE=-O2 -g -fno-omit-frame-pointer " \
    "QMAKE_LFLAGS_RELEASE+=-g"; then
    echo "Release qmake 失败"
    exit 1
fi

if ! make -j$(($(nproc)/2)); then
    echo "Release 编译失败"
    ls -l .
    exit 1
fi
#将所有.o copy回newaoe/release目录
cd ../
mkdir -p newaoe/release
cp newaoe_copy/*.o newaoe/release
rm -f newaoe/release/UsrAI.o
#将moc_*.cpp文件copy到newaoe/release目录下
cp newaoe_copy/moc_*.cpp newaoe/
cp newaoe_copy/ui_*.h newaoe/
# 打印编译结果
echo "Release 版本编译完成"
echo "Release 对象文件：newaoe/release/"
###############################################
# 编译Debug版本
##############################################3
cd newaoe_copy

# 清理旧的编译文件
make clean

# 编译Debug版本
if ! qmake \
    CONFIG+=debug \
    CONFIG-=release \
    "QMAKE_CXXFLAGS_DEBUG=-O0 -g3 -fno-omit-frame-pointer -fno-optimize-sibling-calls -fsanitize=address,undefined -D_GLIBCXX_ASSERTIONS" \
    "QMAKE_LFLAGS_DEBUG+=-fsanitize=address,undefined"; then
    echo "Debug qmake 失败"
    exit 1
fi

if ! make -j$(($(nproc)/2)); then
    echo "Debug 编译失败"
    ls -l .
    exit 1
fi

cd ..

mkdir -p newaoe/debug
cp newaoe_copy/*.o newaoe/debug/
rm -f newaoe/debug/UsrAI.o

# 打印编译结果
echo "Debug 版本编译完成"
echo "Debug 对象文件：  newaoe/debug/"