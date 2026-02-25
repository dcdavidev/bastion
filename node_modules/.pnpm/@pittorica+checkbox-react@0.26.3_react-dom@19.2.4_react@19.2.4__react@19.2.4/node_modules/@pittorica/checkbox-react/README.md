# @pittorica/checkbox-react

The `Checkbox` component.

## Installation

\`\`\`bash
npm install @pittorica/checkbox-react
\`\`\`

## Usage

\`\`\`jsx
import { Checkbox } from '@pittorica/checkbox-react';

function MyForm() {
return (

<form>
<Checkbox
label="Accept Terms and Conditions"
name="terms"
required // Mark as required
color="source" // Default color is now 'source'
onChange={(checked) => console.log('Agreed:', checked)}
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
