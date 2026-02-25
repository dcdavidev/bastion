# @pittorica/code-react

The `Code` component.

## Installation

```bash
npm install @pittorica/code-react
```

## Usage

```jsx
import { Code } from '@pittorica/code-react';

function MyComponent() {
  return (
    <>
      <Code>npm install pittorica</Code>

      <Code language="typescript" filename="Example.ts" showLineNumbers>
        {\`const hello = "world";
console.log(hello);\`}
      </Code>
    </>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
