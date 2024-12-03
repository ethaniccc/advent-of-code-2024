reports = []
max_threshold = 3

def parse_input():
    for line in open("input").read().split("\n"):
        reports.append([int(num) for num in line.split(" ")])

def determine_report_safe(report, tolerate):
    step = 0
    max_index = len(report) - 1

    for index in range(1, len(report)):
        prev, curr = report[index - 1], report[index]
        diff = abs(curr - prev)

        if diff > max_threshold or diff == 0:
            if tolerate:
                return check_report_with_tolerance(report, index)
            return False
        
        if step == 0:
            # We have not yet determined wether this report is incrementing or decrementing.
            next_num = report[index + 1]
            if prev > curr:
                step = -1
                if next_num > curr:
                    if tolerate:
                        return check_report_with_tolerance(report, index)
                    return False
            else:
                step = 1
                if next_num < curr:
                    if tolerate:
                        return check_report_with_tolerance(report, index)
                    return False
        elif step == 1:
            # If the step is positive, we should be expecting an incrementing report.
            if curr < prev:
                if tolerate:
                    return check_report_with_tolerance(report, index)
                return False
        else:
            # If the step is negative we should be expecting a decrementing report.
            if curr > prev:
                if tolerate:
                    return check_report_with_tolerance(report, index)
                return False

    return True


def check_report_with_tolerance(report, index):
    prev_index = index - 1
    r1 = report[:index]
    r1.extend(report[index+1:])
    r2 = report[:prev_index]
    r2.extend(report[prev_index+1:])

    return determine_report_safe(r1, False) or determine_report_safe(r2, False)

def safe_reports(use_problem_dampener):
    safe = 0
    for report in reports:
        if determine_report_safe(report, use_problem_dampener):
            safe += 1
    return safe

def main():
    parse_input()
    print("Safe reports w/o dampener:", safe_reports(False))
    print("Safe reports w/ dampener:", safe_reports(True))

main()