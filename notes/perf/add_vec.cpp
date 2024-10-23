#include <algorithm>
#include <iostream>
#include <vector>

constexpr int vec_size = 10000;

void mul(std::vector<int>& a, std::vector<int>& b, std::vector<int>& c)
{
    std::transform(
        a.begin(),
        a.end(),
        b.begin(),
        c.begin(),
        [](auto op1, auto op2) {return op1 * op2;}
        );
}

std::vector<int> A(vec_size, 1);
std::vector<int> B(vec_size, 2);
std::vector<int> C(vec_size);

int main()
{
    mul(A, B, C);
    for (auto const& v : C)
    {
        if (v != 2)
        {
            std::cout << "Oops, something happened!\n";
            return 1;
        }
    }
    return 0;
}
