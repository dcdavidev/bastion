# @pittorica/switch-react

The `Switch` component.

## Installation

\`\`\`bash
npm install @pittorica/switch-react
\`\`\`

## Usage

\`\`\`jsx
import { Switch } from '@pittorica/switch-react';

function MyForm() {
return (

<form>
<Switch
label="Enable notifications"
name="notifications"
required // Mark as required
color="source" // Default color is now 'source'
/>
<button type="submit">Submit</button>
</form>
);
}
\`\`\`

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
