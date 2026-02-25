# @pittorica/text-field-react

The `TextField` component.

## Installation

\`\`\`bash
npm install @pittorica/text-field-react
\`\`\`

## Usage

\`\`\`jsx
import { TextField } from '@pittorica/text-field-react';

function MyForm() {
return (

<form>
<TextField.Root
label="Username"
name="username"
required // Mark as required
color="source" // Default color is now 'source'
size="sm"
helperText="Enter your username." >
<TextField.Input placeholder="e.g. john_doe" />
</TextField.Root>
<button type="submit">Submit</button>
</form>
);
}
\`\`\`

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
