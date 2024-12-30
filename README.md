# UCT Archive Page

This is the README for the ICS4U SDP. The goal of the project was to develop an archive page for United Coding Tournament (UCT).

Please note that the full website is not open source. Only the files that were strictly created for this project are public. You can find the full website repository at https://github.com/United-Coding-Tournament.

The following files were modified for this project:
* gen.go
* admin.go
* main.go

The following files were created for this project:
* archive.go
* problem.go
* functions.go
* adminarchive.html
* archive.html
* problempage.html

# Features

The user may
* View problems
* Filter problems by
    * Name
    * Status
    * Difficulty
    * Tags
* Submit solutions and receive answer feedback
* Add/delete problems if user has admin permissions

# Installation

* VS Code
    * Navigate to the [Visual Studio Code setup guide](https://code.visualstudio.com/docs/setup/setup-overview) to set up Visual Studio Code
* Postgres 14
    * Navigate to the [Postgres 14 download page](https://www.postgresql.org/download/)
    * Download version 14.8
    * During the installation, make sure of the following:
        * pgAdmin 4 is in the selected components to install
        * Host is **localhost**
        * User Name is **postgres**
        * Port is **Port 5432**
        * Note down the passwords used for superuser and postgres
* Go (lang)
    * Navigate to the [intallation guide](https://go.dev/doc/install) to set up Go

# Usage

Unfortunately, the archive page can no longer be accessed as the UCT website is no longer running.

## Running the Code Locally

Unfortunately this contains some sensitive details that are unable to be shared as the full website is not open source.

# Known Bugs

* File input for uploading new problems does not include error checking for correct file type
* Lack of checking of the content of uploaded files for new problems
* Text input parameters for uploading new problems do not include error checking for specific style guidelines
* No cues on interface if filters are currently active in the archive page apart from the URL
* If admin is filling out the form to add the problem and admin closes the form and filters the problem list, the form will be cleared
* There is no error checking for inserting duplicate problems
    * Database will throw an error, although code does keep running
* Errors relating to database actions are logged but the code does not stop running

# Credits

## Contributers

* Justin Dai
* Charisman Aravinthan

## Sources used for learning
| Source | Purpose |
|-|-|
| Muhammad Wasif Kamran | To learn how to program in Golang and to understand how to use PostgreSQL database with pgAdmin4 |
| Muhammad Wasif Kamran | To learn about how the current UCT website is programmed |
| Muhammad Wasif Kamran | To learn to use git commands |
| “Documentation,” Go, https://go.dev/doc/ (accessed June 1, 2023). | To learn about different Go packages and built-in functions |
| “HTML Tutorial,” w3schools, https://www.w3schools.com/html/default.asp (accessed Jun. 2, 2023). | Used to learn HTML/Tailwind CSS |

## Sources used for code
| Source | Purpose |
|-|-|
| “Tailwind CSS Dropdown,” Flowbite, https://flowbite.com/docs/components/dropdowns/ (accessed May 31, 2023). | Used code for basic structure of dropdown buttons |
| “Tailwind CSS search input,” Flowbite, https://flowbite.com/docs/forms/search-input/ (accessed May 31, 2023). | Used code for basic structure of search input |
| “Tailwind CSS table,” Flowbite, https://flowbite.com/docs/components/tables/ (accessed May 31, 2023). | Used code for basic structure of tables |
| UCT, https://theuct.ca (accessed June 2, 2023). | Used code in other files of the website for code such as the problem submission form |