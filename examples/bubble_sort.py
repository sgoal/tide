def bubble_sort(arr):
    n = len(arr)
    for i in range(n):
        for j in range(0, n - i - 1):
            if arr[j] > arr[j + 1]:
                arr[j], arr[j + 1] = arr[j + 1], arr[j]


def quick_sort(arr):
    if len(arr) <= 1:
        return arr
    pivot = arr[len(arr) // 2]
    left = [x for x in arr if x < pivot]
    middle = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]
    return quick_sort(left) + middle + quick_sort(right)




def swap_sort(arr):
    n = len(arr)
    for i in range(n):
        for j in range(i + 1, n):
            if arr[i] > arr[j]:
                arr[i], arr[j] = arr[j], arr[i]


if __name__ == '__main__':
    test_arr_bubble = [64, 34, 25, 12, 22, 11, 90]
    print("Bubble sort result:", end=" ")
    bubble_sort(test_arr_bubble)
    print(test_arr_bubble)

    test_arr_quick = [64, 34, 25, 12, 22, 11, 90]
    print("Quick sort result:", end=" ")
    sorted_arr = quick_sort(test_arr_quick)
    print(sorted_arr)

    test_arr_swap = [64, 34, 25, 12, 22, 11, 90]
    print("Swap sort result:", end=" ")