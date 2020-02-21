Maximize total score of all scanned books. 
- Each library needs to be signed up before it can ship books

# Books

_B_ books with ID's 0 -> _B_-1. Multiple copies of the same book exist, we only need one scanned. 
- Books scanned must be unique. IE points only awarded for each unique books. 
# Libraries

_L_ librareis from 0 -> _L_-1. 

Attributes
    - set of **Books**
    - **time** in days so sign the library up
    - **Number** of **Books** that can be scanned each day
  

# Time 
0 to D-1 : can ship books. 
First library signup on day 0


# Library signup (Basically a mutex lock)

- Each library must be signed up.
- Only **one** may be signed up at a time
- Libraries can be signed up in any order
- On the first day after signed up **Books** can be **scanned**

# Scanning

- Scanning is measured in 1 day intervals. 
- The Scanning facility can scan an infinite number of books in a day
- **Bottleneck is the libraries**



# Example

6 2 7 - Books, Libraries, Days
1 2 3 6 5 4 - Score of each respective book

->lib 0
5 2 2 - number of books, number of days for signup, scannable books per day
0 1 2 3 4 - book ids

->lib 1
4 3 1 - number of books, number of days for signup, scannable books per day
0 2 3 5 - book ids


# In short

X days, each library needs to signup before we can start shipping books.

Signup 1 library at a time, scanning books happens concurrently & each library has a max quota.