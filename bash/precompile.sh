 #先copy一下newaoe目录
cp -r newaoe newaoe_copy
#修复所有文件的大小写
fixcase -f ./newaoe_copy
#切换到newaoe_copy目录下
cd newaoe_copy
#添加qt5.9.2的bin目录到PATH
export PATH="$PATH:/opt/qt5.9.2/bin/"
#编译
if ! qmake QMAKE_CXXFLAGS+=" -O2 -g -fno-omit-frame-pointer" QMAKE_LFLAGS+=" -g" || ! make -j$(($(nproc)/2)); then
    echo "编译失败"
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