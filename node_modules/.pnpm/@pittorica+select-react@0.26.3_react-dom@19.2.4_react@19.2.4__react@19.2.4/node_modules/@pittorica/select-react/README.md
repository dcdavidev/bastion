# @pittorica/select-react

The `Select` component.

## Installation

\`\`\`bash
npm install @pittorica/select-react
\`\`\`

## Usage

\`\`\`jsx
import { Select } from '@pittorica/select-react';

function MyForm() {
return (

<form>
<Select.Root
label="Country"
name="country"
required // Mark as required
color="source" // Default color is now 'source'
size="sm"
defaultValue="us" >
<Select.Content>
<option value="">Select a country</option>
<option value="us">United States</option>
<option value="ca">Canada</option>
<option value="mx">Mexico</option>
</Select.Content>
</Select.Root>
<button type="submit">Submit</button>
</form>
);
}
\`\`\`

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
