# PDF Downloader (Uge 5 Projekt)

The PDF Downloader task from week 5.

Works the following way:
- Takes a specific Excel speadsheet (provided in the data folder) as input. 
- Then reads each row as a "report" with data from relevant columns.
- Then downloads all reports in parallel with a helpful progress bar for each download.
- Then writes the result of each download to a metadata.xlsx in the output dir

## Building

Simply use the command `go build`

## Usage

The following commandline arguments are required for the program to work.

- **input_data**=_excel_spreadsheet_path_
- **output_dir**=_output_directory_

Note:  
If using VS Code, you can also just launch it in the debugger, which has the arguments supplied.
