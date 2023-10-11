#include<>
//快排   
//生成随机数  
#include <iostream>
void creatRandomNumber (int *a, int n) {
    srand((unsigned)time(NULL));
    for (int i = 0; i < n; i++) {
        a[i] = rand() % 100;
    }
}
//快速排序  
void quickSort(int *a, int left, int right) {
    if (left >= right) {
        return;
    }
    int i = left;
    int j = right;
    int key = a[left];
    while (i < j) {
        while (i < j && key <= a[j]) {
            j--;
        }
        a[i] = a[j];
        while (i < j && key >= a[i]) {
            i++;
        }
        a[j] = a[i];
    }
    a[i] = key;
    quickSort(a, left, i - 1);
    quickSort(a, i + 1, right);
}
