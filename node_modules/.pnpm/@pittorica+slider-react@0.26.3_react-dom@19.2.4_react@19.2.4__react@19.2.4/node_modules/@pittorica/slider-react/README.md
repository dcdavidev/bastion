# @pittorica/slider-react

The `Slider` component.

## Installation

\`\`\`bash
npm install @pittorica/slider-react
\`\`\`

## Usage

\`\`\`jsx
import { Slider } from '@pittorica/slider-react';

function MyForm() {
return (

<form>
<Slider
min={0}
max={100}
defaultValue={50}
color="source" // Default color is now 'source'
required // Mark as required
label="Volume"
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
