package main

import (
    "fmt"

    "github.com/valkey-io/valkey-glide/go/api"
)

func main() {
    host := "clustercfg.memorydb-singleregion.nkdf1r.memorydb.us-east-1.amazonaws.com"
    port := 6379

    config := api.NewGlideClusterClientConfiguration().
        WithAddress(&api.NodeAddress{Host: host, Port: port}).
        WithUseTLS(true)

    client, err := api.NewGlideClusterClient(config)
    if err != nil {
        fmt.Printf("Connection error: %v\n", err)
        return
    }
    defer client.Close()

    // Test connection
    res, err := client.Ping()
    if err != nil {
        fmt.Printf("Ping error: %v\n", err)
        return
    }
    fmt.Printf("Ping response: %s\n", res)

    // Set a value
    setRes, err := client.Set("test-key", "hello world")
    if err != nil {
        fmt.Printf("Set error: %v\n", err)
        return
    }
    fmt.Printf("Set response: %v\n", setRes)

    // Get a value
    value, err := client.Get("test-key")
    if err != nil {
        fmt.Printf("Get error: %v\n", err)
        return
    }
    fmt.Printf("Retrieved value: %s\n", value)

    // Delete a key (Del takes a slice of strings)
    delRes, err := client.Del([]string{"test-key"})
    if err != nil {
        fmt.Printf("Delete error: %v\n", err)
        return
    }
    fmt.Printf("Delete response: %v\n", delRes)

    // First, let's set multiple keys
    keys := []string{"key1", "key2", "key3"}
    for i, key := range keys {
        _, err := client.Set(key, fmt.Sprintf("value%d", i+1))
        if err != nil {
            fmt.Printf("Error setting %s: %v\n", key, err)
            return
        }
    }

    // Now get multiple keys at once
    values, err := client.MGet(keys)
    if err != nil {
        fmt.Printf("Batch get error: %v\n", err)
        return
    }

    // Process results
    for i, value := range values {
        fmt.Printf("Key: %s, Value: %v\n", keys[i], value)
    }



    // RPUSH - Append to the right (end) of the list

    listKey := "mylist"

    // Single value append
    count, err := client.RPush(listKey, []string{"value1"})
    if err != nil {
        fmt.Printf("Error appending to list: %v\n", err)
        return
    }
    fmt.Printf("RPush result (list length): %d\n", count)

    // Multiple values append
    moreValues := []string{"value2", "value3"}
    count, err = client.RPush(listKey, moreValues)
    if err != nil {
        fmt.Printf("Error appending multiple values: %v\n", err)
        return
    }
    fmt.Printf("RPush multiple result (list length): %d\n", count)

    // Get all values from the list
    listValues, err := client.LRange(listKey, 0, -1)
    if err != nil {
        fmt.Printf("Error getting list values: %v\n", err)
        return
    }

    // Print results
    fmt.Println("List values:")
    for _, val := range listValues {
        fmt.Printf("%v\n", val)
    }

    // Get list length
    length, err := client.LLen(listKey)
    if err != nil {
        fmt.Printf("Error getting list length: %v\n", err)
        return
    }
    fmt.Printf("List length: %d\n", length)

    // Optional: Delete the list when done
    _, err = client.Del([]string{listKey})
    if err != nil {
        fmt.Printf("Error deleting list: %v\n", err)
        return
    }

   //LTRIM

    // Add test data using direct string slice
    // Declare all variables upfront
    var (
        trimKey     = "trimlist"
        trimResult  string
    )

    count, err = client.RPush(trimKey, []string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10"})
    if err != nil {
        fmt.Printf("Error adding values: %v\n", err)
        return
    }
    fmt.Printf("Added %d values\n", count)

    // Print initial list
    listValues, err = client.LRange(trimKey, 0, -1)
    if err != nil {
        fmt.Printf("Error getting initial list: %v\n", err)
        return
    }
    fmt.Println("Initial list:", listValues)

    // Example 1: Remove first 2 elements
    trimResult, err = client.LTrim(trimKey, 2, -1)
    if err != nil {
        fmt.Printf("Error trimming first elements: %v\n", err)
        return
    }
    fmt.Printf("Trim result 1: %v\n", trimResult)

    listValues, err = client.LRange(trimKey, 0, -1)
    if err != nil {
        fmt.Printf("Error getting list after first trim: %v\n", err)
        return
    }
    fmt.Println("After removing first 2:", listValues)

    // Example 2: Remove last 2 elements
    length, err = client.LLen(trimKey)
    if err != nil {
        fmt.Printf("Error getting length: %v\n", err)
        return
    }

    trimResult, err = client.LTrim(trimKey, 0, length-3)
    if err != nil {
        fmt.Printf("Error trimming last elements: %v\n", err)
        return
    }
    fmt.Printf("Trim result 2: %v\n", trimResult)

    listValues, err = client.LRange(trimKey, 0, -1)
    if err != nil {
        fmt.Printf("Error getting list after second trim: %v\n", err)
        return
    }
    fmt.Println("After removing last 2:", listValues)

    // Example 3: Keep only middle section
    trimResult, err = client.LTrim(trimKey, 1, 3)
    if err != nil {
        fmt.Printf("Error trimming to middle section: %v\n", err)
        return
    }
    fmt.Printf("Trim result 3: %v\n", trimResult)

    listValues, err = client.LRange(trimKey, 0, -1)
    if err != nil {
        fmt.Printf("Error getting list after middle trim: %v\n", err)
        return
    }
    fmt.Println("After keeping only middle section:", listValues)

    // Example 4: Keep only last 2 elements
    trimResult, err = client.LTrim(trimKey, -2, -1)
    if err != nil {
        fmt.Printf("Error trimming to last elements: %v\n", err)
        return
    }
    fmt.Printf("Trim result 4: %v\n", trimResult)

    listValues, err = client.LRange(trimKey, 0, -1)
    if err != nil {
        fmt.Printf("Error getting final list: %v\n", err)
        return
    }
    fmt.Println("After keeping only last 2:", listValues)

    //HSET

   var (
        hashKey = "user:123"
        result  int64
    )

    // Set single field in hash using map
    singleField := map[string]string{
        "name": "John",
    }
    result, err = client.HSet(hashKey, singleField)
    if err != nil {
        fmt.Printf("Error setting single field: %v\n", err)
        return
    }
    fmt.Printf("Single field set result: %v\n", result)

    // Set multiple fields at once
    fields := map[string]string{
        "age":    "30",
        "city":   "New York",
        "email":  "john@example.com",
    }

    // Set multiple fields
    result, err = client.HSet(hashKey, fields)
    if err != nil {
        fmt.Printf("Error setting multiple fields: %v\n", err)
        return
    }
    fmt.Printf("Multiple fields set result: %v\n", result)

    // Get all fields from hash
    allFields, err := client.HGetAll(hashKey)
    if err != nil {
        fmt.Printf("Error getting all fields: %v\n", err)
        return
    }
    fmt.Println("All hash fields:", allFields)

    // Get specific fields
    name, err := client.HGet(hashKey, "name")
    if err != nil {
        fmt.Printf("Error getting name: %v\n", err)
        return
    }
    fmt.Printf("Name: %v\n", name)

    // Check if field exists
    exists, err := client.HExists(hashKey, "age")
    if err != nil {
        fmt.Printf("Error checking field existence: %v\n", err)
        return
    }
    fmt.Printf("Age field exists: %v\n", exists)

    // Get number of fields in hash
    length, err = client.HLen(hashKey)

    if err != nil {
        fmt.Printf("Error getting hash length: %v\n", err)
        return
    }
    fmt.Printf("Number of fields: %d\n", length)
  //  fmt.Printf("HLen method type: %T\n", client.HLen)

    // Delete specific fields
    result, err = client.HDel(hashKey, []string{"email"})
    if err != nil {
        fmt.Printf("Error deleting field: %v\n", err)
        return
    }
    fmt.Printf("Delete result: %v\n", result)

    // Get all fields after deletion
    allFields, err = client.HGetAll(hashKey)
    if err != nil {
        fmt.Printf("Error getting final fields: %v\n", err)
        return
    }
    fmt.Println("Hash fields after deletion:", allFields)

    fmt.Println("Operations completed successfully!")
}
