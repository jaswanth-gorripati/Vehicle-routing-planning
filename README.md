# Vehicle routing planning

This project is a Vehicle routing planning for drivers using a heuristic approach. It reads input from a text file, assigns loads to drivers while minimizing the total cost, and prints the optimized load assignment for each driver.

## Features

- Parses input from a text file.
- Calculates distances using the Euclidean distance formula.
- Assigns loads to drivers based on distance constraints.
- Uses a heuristic optimization approach to minimize costs.
- Prints the optimized solution.

## Getting Started

### Prerequisites

- Go programming language installed on your machine. You can download and install it from [here](https://golang.org/dl/).

### Installation

1. Clone this repository:

   ```sh
   git clone https://github.com/jaswanth-gorripati/Vehicle-routing-planning.git
   ```

2. Navigate to the project directory:

   ```sh
   cd Vehicle-routing-planning
   ```

## Usage

### Input File Format

The input file should be a text file with tab-separated values and the following structure:

```txt
loadNumber pickup dropoff
1 (0.0,1.0) (2.0,3.0)
2 (1.0,1.0) (3.0,3.0)
```

- `loadNumber`: Unique identifier for the load.
- `pickup`: Coordinates of the pickup location.
- `dropoff`: Coordinates of the dropoff location.

### Running the Program

1. Ensure your input file is correctly formatted as described above.

2. Run the program with the path to your input file as an argument:

   ```sh
   go run main.go path/to/your/inputfile.txt
   ```

   Replace `path/to/your/inputfile.txt` with the actual path to your input file.

### Example

1. Create an input file `loads.txt`:

    ```txt
    loadNumber pickup       dropoff
    1          (0.0,1.0)    (2.0,3.0)
    2          (1.0,1.0)    (3.0,3.0)
    3          (2.0,2.0)    (4.0,5.0)
    ```

2. Run the program:

    ```sh
    go run main.go loads.txt
    ```

3. The program will print the optimized load assignments for each driver.

### Sample Output

```txt
 [1,2]
 [3]
```
This indicates that the first driver is assigned loads 1 and 2, while the second driver is assigned load 3.
