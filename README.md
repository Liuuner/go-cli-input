# Go Cli Input

### TBD
- [x] Select
- [x] Text
- [x] Text Sensitive
- [x] Checkboxes
- [x] Boolean [Y/n] [y/N] [y/n] ...

#### More Ideas
- [ ] Dropdown Menu
- [ ] Number Input
- [ ] Date Picker
- [ ] Time Picker
- [ ] File Upload
- [ ] Slider
- [ ] Text Area
- [ ] Rating
- [ ] Search Box
- [ ] Color Picker
- [ ] Tags

"Tags" refer to a method of categorizing or labeling items, allowing users to add keywords or phrases that describe the content or context of a particular item. In the context of input types, a **Tags input** typically allows users to enter multiple tags or keywords associated with an item, often separated by commas or spaces.

### Features of Tags Input:
- **Multi-Selection**: Users can input multiple tags, which may be displayed as individual tags in a list or cloud format.
- **Autocomplete**: As users type, suggestions may appear based on previously used tags or a predefined list, helping to speed up the tagging process.
- **Removable**: Users can usually remove tags by clicking on a "remove" icon or pressing the backspace key.
- **Validation**: You can implement checks to ensure that tags meet certain criteria (e.g., length, uniqueness).

### Example Use Cases:
- **Blog Posts**: Allowing authors to tag their posts with relevant keywords (e.g., "Go", "Programming", "Web Development").
- **Project Management**: Users can tag tasks or projects for better organization and filtering (e.g., "urgent", "in-progress", "completed").
- **Social Media**: Users can tag content or friends in posts (e.g., "vacation", "family", "food").

### Example UI:
In a tags input field, users might see something like this after entering their tags:

```
Tags: [Go] [Programming] [Web Development] [X] [Remove]
```

```
// Like list of inputs?
Tags: 
    - tag1
    - tag2
    - tag3
```
