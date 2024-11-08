# PDF Downloader (Uge 5 Projekt)

The PDF Download task from week 5.
Takes a specific Excel sheet (in the data folder) as input, reads each row as a "report" with relevant data.
Then downloads all reports in parallel with a helpful progress bar for each download.

## Building

Simply use the command `go build`

## Usage

The following commandline arguments are required for the program to work.

- **input_data**=_excel_spreadsheet_path_
- **output_dir**=_output_directory_

Note:  
If using VS Code, you can also just launch it in the debugger, which has the arguments supplied.
