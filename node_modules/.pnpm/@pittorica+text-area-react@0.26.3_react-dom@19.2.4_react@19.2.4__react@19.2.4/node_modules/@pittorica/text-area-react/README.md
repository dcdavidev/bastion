# @pittorica/text-area-react

The `TextArea` component.

## Installation

\`\`\`bash
npm install @pittorica/text-area-react
\`\`\`

## Usage

\`\`\`jsx
import { TextArea } from '@pittorica/text-area-react';

function MyForm() {
return (

<form>
<TextArea.Root
label="Comments"
name="comments"
required // Mark as required
color="source" // Default color is now 'source'
size="sm"
helperText="Enter your feedback below." >
<TextArea.Content placeholder="Type your message here..." />
</TextArea.Root>
<button type="submit">Submit</button>
</form>
);
}
\`\`\`

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
