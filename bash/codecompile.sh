export PATH="$PATH:/opt/qt5.9.2/bin/" 
export QT="/opt/qt5.9.2" 
export QTINCLUDE="/opt/qt5.9.2/include" 
export DebugMode={DebugMode}
#先copy一下newaoe_copy目录
mkdir -p newaoe_copy
cd project
find ./ \
    \( -name "*.h" -o -name "*.cpp" -o -name "*.hpp" \) \
    -exec cp --parents {{}} ../newaoe_copy \;
cd ../
#把UsrAI.h和UsrAI.cpp copy到newaoe_copy目录下
cp build/UsrAI.h build/UsrAI.cpp newaoe_copy/
#切换到newaoe_copy目录下
cd newaoe_copy
#修复大小写问题
fixcase -f ./ >/dev/null
# 编译用户代码
if [ $DebugMode = true ]; then
    g++ -c UsrAI.cpp \
    -Og \
    -g3 \
    -fno-omit-frame-pointer \
    -fno-optimize-sibling-calls \
    -fPIC \
    -fsanitize=address,undefined \
    -Wall -Wextra -Wpedantic \
    -D_GLIBCXX_ASSERTIONS \
    -I./ \
    -I${{QTINCLUDE}} \
    -I${{QTINCLUDE}}/QtCore \
    -I${{QTINCLUDE}}/QtMultimedia \
    -I${{QTINCLUDE}}/QtWidgets \
    -I${{QTINCLUDE}}/QtGui \
    -I${{QTINCLUDE}}/QtNetwork &&

    # 链接公共 .o 文件
    g++ UsrAI.o ../project/debug/*.o \
        -Og \
        -g3 \
        -fno-omit-frame-pointer \
        -fsanitize=address,undefined \
        -o newAOE \
        -L/opt/qt5.9.2/lib \
        -lQt5Widgets \
        -lQt5Gui \
        -lQt5Core \
        -lQt5Multimedia \
        -lQt5Network
else
    g++ -c UsrAI.cpp \
    -O2 \
    -fPIC \
    -g \
    -fno-omit-frame-pointer \
    -I./ \
    -I${{QTINCLUDE}} \
    -I${{QTINCLUDE}}/QtCore \
    -I${{QTINCLUDE}}/QtMultimedia \
    -I${{QTINCLUDE}}/QtWidgets \
    -I${{QTINCLUDE}}/QtGui \
    -I${{QTINCLUDE}}/QtNetwork &&

    # 链接公共 .o 文件
    g++ UsrAI.o ../project/release/*.o \
        -O2 \
        -g \
        -fno-omit-frame-pointer \
        -o newAOE \
        -L/opt/qt5.9.2/lib \
        -lQt5Widgets \
        -lQt5Gui \
        -lQt5Core \
        -lQt5Multimedia \
        -lQt5Network
fi

#把newAOE copy回buildDir目录下
cp newAOE ../build/  >/dev/null