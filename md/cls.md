#### 基本数据类型

```cpp
#include <iostream>
using namespace std;

int main()
{
    cout << "sizeof(bool)          : " << sizeof(bool) << endl;
    cout << "sizeof(char)          : " << sizeof(char) << endl;
    cout << "sizeof(int)           : " << sizeof(int) << endl;
    cout << "sizeof(unsigned int)  : " << sizeof(unsigned int) << endl;
    cout << "sizeof(short int)     : " << sizeof(short int) << endl;
    cout << "sizeof(long int)      : " << sizeof(long int) << endl;
    cout << "sizeof(float)         : " << sizeof(float) << endl;
    cout << "sizeof(double)        : " << sizeof(double) << endl;

    cout << "min(bool)          : " << numeric_limits<bool>::min() << endl;
    cout << "min(int)           : " << numeric_limits<int>::min() << endl;
    cout << "min(unsigned int)  : " << numeric_limits<unsigned int>::min() << endl;
    cout << "min(short int)     : " << numeric_limits<short int>::min() << endl;
    cout << "min(long int)      : " << numeric_limits<long int>::min() << endl;
    cout << "min(float)         : " << numeric_limits<float>::min() << endl;
    cout << "min(double)        : " << numeric_limits<double>::min() << endl;

    cout << "max(bool)          : " << numeric_limits<bool>::max() << endl;
    cout << "max(int)           : " << numeric_limits<int>::max() << endl;
    cout << "max(unsigned int)  : " << numeric_limits<unsigned int>::max() << endl;
    cout << "max(short int)     : " << numeric_limits<short int>::max() << endl;
    cout << "max(long int)      : " << numeric_limits<long int>::max() << endl;
    cout << "max(float)         : " << numeric_limits<float>::max() << endl;
    cout << "max(double)        : " << numeric_limits<double>::max() << endl;

    return 0;
}

```

输出结果

```bash
sizeof(bool)          : 1
sizeof(char)          : 1
sizeof(int)           : 4
sizeof(unsigned int)  : 4
sizeof(short int)     : 2
sizeof(long int)      : 8
sizeof(float)         : 4
sizeof(double)        : 8
min(bool)          : 0
min(int)           : -2147483648
min(unsigned int)  : 0
min(short int)     : -32768
min(long int)      : -9223372036854775808
min(float)         : 1.17549e-38
min(double)        : 2.22507e-308
max(bool)          : 1
max(int)           : 2147483647
max(unsigned int)  : 4294967295
max(short int)     : 32767
max(long int)      : 9223372036854775807
max(float)         : 3.40282e+38
max(double)        : 1.79769e+308
```



#### 指针

##### 分类

- 