# @pittorica/react

This package is a collection of all the React components in the Pittorica design system.

## Installation

To install the entire component library, run the following command:

```bash
npm install @pittorica/react
```

You will also need to install the core `pittorica` package which contains the CSS.

```bash
npm install pittorica
```

## Usage

Once installed, you can import the components from `@pittorica/react`.

```jsx
import 'pittorica/reset';
import 'pittorica/tokens';
import { Button, Card } from '@pittorica/react';

function MyApp() {
  return (
    <Card>
      <Button>Click me</Button>
    </Card>
  );
}
```

## Available Components

This package includes the following components:

- `AlertDialog`
- `AspectRatio`
- `Avatar`
- `Badge`
- `Blockquote`
- `Box`
- `Button`
- `Callout`
- `Card`
- `Carousel`
- `Checkbox`
- `CheckboxCard`
- `CheckboxGroup`
- `Chip`
- `Code`
- `Container`
- `ContextMenu`
- `DataList`
- `Dialog`
- `Divider`
- `DropdownMenu`
- `Em`
- `Flex`
- `Grid`
- `Heading`
- `HoverCard`
- `IconButton`
- `Inset`
- `Kbd`
- `Link`
- `Popover`
- `Progress`
- `Quote`
- `Radio`
- `RadioCard`
- `RadioGroup`
- `Section`
- `SegmentedControl`
- `Select`
- `Sheet`
- `Skeleton`
- `Slider`
- `Stack`
- `Strong`
- `Switch`
- `Table`
- `Tabs`
- `Text`
- `TextArea`
- `TextField`
- `Theme`
- `Toast`
- `Tooltip`

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
