# DG PBT Parser

## Description

The intention of this project is to create an executable that can synthesise info from PBT spreadsheets (portal runsheet.xlsx and weekly invoices.csv) into the required format for Divers Accounting.

As it currently stands, making spreadsheets each month is an extremely tedious job that honestly nobody sane should be doing. Unfortunately, there is no public PBT API (or maybe we just don't have access to it (¬_¬")) which limits everything to just the spreadsheets.

Some of the data is a bit confidential, so all the non *_test.json files in the config folder are hidden. I wrote this in my own time, so it probably won't be written amazingly.

## TODO

- Make it easier to edit the data in the database? (such as correct unknown sortby codes)
- Iterate through all "runsheet_exporting.xlsx" files in a folder
- Use invoice spreadsheet files to fill in price columns
- Iterate through all invoice spreadsheet files in a folder
- Export final data in .xlsx format with all the correct formula
- Create the other two PBT accounts (repeating steps 2-5?)
- GUI?