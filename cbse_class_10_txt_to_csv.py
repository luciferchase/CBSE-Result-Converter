# A luciferchase original

import csv

# Edit this to the path to your result file
path = ""

with open(path, "r") as f:
    raw_data = f.read()
    raw_data = raw_data.split()

parsed_data = [["Roll", "Gender", "Name", "English (184)", "Hindi (002)", "Maths (041)", 
"Science (086)", "Social Sc. (087)", "IT (402)", "Total", "Percentage", "Result"]]

for i in range(len(raw_data)):
    # Filter anything which looks like a roll number
    if (raw_data[i][:4] == "1413"):
        # Parse individual student's data

        roll, gender = (raw_data[i], raw_data[i + 1])
        
        # If a student has only one word name
        if (raw_data[i + 3] == "184"):
            name = raw_data[i + 2]
            i -= 1
        # If a student has three word name
        elif (raw_data[i + 5] == "184"):
            name = " ".join(raw_data[i + 2: i + 5])
            i += 1
        # Rare but if a student has four word name
        elif (raw_data[i + 6] == "184"):
            name = " ".join(raw_data[i + 2: i + 6])
            i += 2
        # Default two word names
        else:
            name = " ".join(raw_data[i + 2: i + 4])
        
        result = raw_data[i + 10]
        if (result not in ["PASS", "COMP"]):
            continue

        marks = [raw_data[i + 11], raw_data[i + 13], raw_data[i + 15], raw_data[i + 17], 
        raw_data[i + 19], raw_data[i + 21]]
        int_marks = [int(i) for i in marks]

        total = str(sum([int(i) for i in marks]))
        percentage = str((int(total) - int(raw_data[i + 21])) / 5)[:4]

        student_data = [roll, gender, name, marks[0], marks[1], marks[2], marks[3], marks[4], 
        marks[5], total, percentage, result]
        parsed_data.append(student_data)

with open("class_10th_result.csv", "w", newline = "\n") as f:
    writer = csv.writer(f)
    writer.writerows(parsed_data)