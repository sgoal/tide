def bubble_sort(arr):
    n = len(arr)
    for i in range(n):
        for j in range(0, n-i-1):
            if arr[j] > arr[j+1]:
                arr[j], arr[j+1] = arr[j+1], arr[j]
    return arr

# 示例
if __name__ == "__main__":
    array = [64, 34, 25, 12, 22, 11, 90]
    print("排序前:", array)
    sorted_array = bubble_sort(array)
    print("排序后:", sorted_array)