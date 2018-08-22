# Intro
This is a go program to fix Cambridgesoft Inventory 13 Bug when adding substance.
The bug will lead to error when adding.Or Cannot display the correct structure after adding.

# The Method
This Program First Search the incorrect data in dataspace then find the SMILES from the database on the Internet.

Then convert SMILES to BASE64 CDX format by using C# API

After that  correct the data in the database

If cannot get any info from internet the program will ask you input manually.

# Useage Tips

1. Open MSSQL sa account and ip connect
2. Add substance using Inventory
3. Run this program

# Speacial Thanks

The lovely friends @ School of Chemistry and Chemical Engineering,Nanjing University.

National Center for Biotechnology Information, U.S. National Library of Medicine