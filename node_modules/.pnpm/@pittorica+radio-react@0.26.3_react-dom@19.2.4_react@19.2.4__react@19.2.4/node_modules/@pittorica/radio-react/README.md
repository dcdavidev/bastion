# @pittorica/radio-react

The `Radio` component.

## Installation

\`\`\`bash
npm install @pittorica/radio-react
\`\`\`

## Usage

\`\`\`jsx
import { Radio } from '@pittorica/radio-react';

function MyForm() {
const [selectedValue, setSelectedValue] = useState('option1');

return (

<form>
<Radio
label="Option 1"
value="option1"
name="myRadioGroup"
checked={selectedValue === 'option1'}
onCheckedChange={(checked) => {
if (checked) setSelectedValue('option1');
}}
required // Mark as required
color="source" // Default color is now 'source'
/>
<Radio
label="Option 2"
value="option2"
name="myRadioGroup"
checked={selectedValue === 'option2'}
onCheckedChange={(checked) => {
if (checked) setSelectedValue('option2');
}}
/>
<button type="submit" disabled={!selectedValue}>Submit</button>
</form>
);
}
\`\`\`

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
