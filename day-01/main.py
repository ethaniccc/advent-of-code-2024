left_nums = []
right_nums = []

def parse_data():
    file_dat = open("input.txt")
    for line in file_dat.read().split("\n"):
        split = line.split(" " * 3)
        left_nums.append(int(split[0]))
        right_nums.append(int(split[1]))

    # Sort out both lists in the parse function, in hindsight, I should've
    # done this with the golang solution as well....
    left_nums.sort()
    right_nums.sort()

def calculate_distance():
    # Both lists should contain the same amount of numbers, given the input file.
    assert len(left_nums) == len(right_nums)

    index = 0
    max_index = len(left_nums) - 1
    total_distance = 0

    while index <= max_index:
        left_num = left_nums[index]
        right_num = right_nums[index]

        total_distance += abs(left_num - right_num)
        index += 1

    return total_distance

def calculate_similarity():
    apperances = {}
    similarity_score = 0

    for num in right_nums:
        # Add the number to the dictionary if it is not already in.
        apperances.setdefault(num, 0)
        apperances[num] += 1

    for num in left_nums:
        # Add the number to the dictionary if it is not already in. This can happen if
        # the number in the left nums list does not exist in the right nums list.
        apperances.setdefault(num, 0)
        similarity_score += apperances[num] * num
    
    return similarity_score

parse_data()
print(f"distance: {calculate_distance()}")
print(f"similarity: {calculate_similarity()}")