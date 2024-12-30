// File: archive.go
// Authors: Charisman Aravinthan and Justin Dai
// Last Modified: June 2, 2023
// Version: 1.1.0
// This file contains additional functions to be used

// Additional functions.
package backend

// Checks if a string is present in an array of strings.
// Parameters: s (slice): array of strings; b (string): element to be checked if present in array;
// Returns: bool: true if b is in s, false otherwise;
// Contributers: Justin Dai.
func contains(s []string, b string) bool {
    for _, a := range s {
        if a == b {
            return true
        }
    }
    return false
}

// Locates the earliest index of a particular string in an array of strings.
// Parameters: s (slice): array of strings; b (string): element to be searched for index;
// Returns: i (int) index of the array, -1 if element is not found;
// Contributers: Justin Dai.
func index(s []string, b string) int {
    for i, a := range s {
        if a == b {
            return i
        }
    }
    return -1
}

// Removes an element from an array of strings. Note that order is not preserved.
// Parameters: s (slice): array of strings; i (index): index of element to be removed;
// Returns: s (slice): updated array with removed element;
// Contributers: Justin Dai.
func remove(s []string, i int) []string {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}