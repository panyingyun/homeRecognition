md .\build_release
pushd .\build_release
del /a CMakeCache.txt
cmake -G "Visual Studio 15 2017 Win64" ^
-DJPEG_INCLUDE_DIR=../dlib/external/libjpeg ^
-DJPEG_LIBRARY=../dlib/external/libjpeg ^
-DPNG_PNG_INCLUDE_DIR=../dlib/external/libpng ^
-DPNG_LIBRARY_RELEASE=../dlib/external/libpng ^
-DZLIB_INCLUDE_DIR=../dlib/external/zlib ^
-DZLIB_LIBRARY_RELEASE=../dlib/external/zlib ^
-DBUILD_SHARED_LIBS:BOOL=FALSE  ^
-DCMAKE_INSTALL_PREFIX="E:/dlib" ^
-DBUILD_TYPE=Release ..

cmake --build . --config Release --target INSTALL
popd
pause