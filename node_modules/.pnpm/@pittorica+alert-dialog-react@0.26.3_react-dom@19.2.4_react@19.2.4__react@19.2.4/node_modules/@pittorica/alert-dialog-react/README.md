# @pittorica/alert-dialog-react

The `AlertDialog` component is used to interrupt the user with a mandatory confirmation or action. It is designed to be accessible and fully themed for both light and dark modes.

## Installation

```bash
npm install @pittorica/alert-dialog-react
```

You will also need to install the core `pittorica` package which contains the CSS.

```bash
npm install pittorica @pittorica/box-react @pittorica/button-react
```

## Usage

```jsx
import { useState } from 'react';
import {
  AlertDialog,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogActions,
} from '@pittorica/alert-dialog-react';
import { Button } from '@pittorica/button-react';
import 'pittorica';

function DeleteAccount() {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button variant="tonal" color="red" onClick={() => setOpen(true)}>
        Delete Account
      </Button>

      <AlertDialog open={open} onClose={() => setOpen(false)}>
        <AlertDialogTitle color="red">
          Are you absolutely sure?
        </AlertDialogTitle>
        <AlertDialogDescription>
          This action cannot be undone. This will permanently delete your
          account and remove your data from our servers.
        </AlertDialogDescription>
        <AlertDialogActions>
          <Button variant="text" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button variant="filled" color="red" onClick={() => setOpen(false)}>
            Yes, delete account
          </Button>
        </AlertDialogActions>
      </AlertDialog>
    </>
  );
}
```

## API Reference

### AlertDialog

The main container for the alert dialog.

| Prop                  | Type         | Default | Description                                      |
| :-------------------- | :----------- | :------ | :----------------------------------------------- |
| `open`                | `boolean`    | `false` | If `true`, the dialog is open.                   |
| `onClose`             | `() => void` | -       | Callback fired when the dialog is closed.        |
| `closeOnOverlayClick` | `boolean`    | `false` | If `true`, closing is allowed via overlay click. |
| `closeOnEsc`          | `boolean`    | `false` | If `true`, closing is allowed via Esc key.       |

### AlertDialogTitle

The title of the alert dialog.

| Prop    | Type     | Default   | Description                                             |
| :------ | :------- | :-------- | :------------------------------------------------------ |
| `color` | `string` | `inherit` | The semantic color of the title (e.g., "red", "amber"). |

### AlertDialogDescription

The description or body text of the alert dialog.

### AlertDialogActions

A container for the action buttons at the bottom of the dialog.

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
