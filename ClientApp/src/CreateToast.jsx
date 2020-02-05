import { SnackbarProvider } from "notistack";
import React from 'react';
import Snackbar from '@material-ui/core/Snackbar';
import MuiAlert from '@material-ui/lab/Alert';

function Alert(props) {
    return <MuiAlert elevation={0} variant="filled" {...props} />
}

export default function CreateToast(props) {
    const [open, setOpen] = [props.open, props.setOpen];
    const path = props.path

    const handleClose = (event, reason) => {
        if (reason === 'clickaway') {
          return;
        }
    
        setOpen(false);
      };

    return (
        <SnackbarProvider>
            <Snackbar open={open} autoHideDuration={4000} onClose={handleClose}>
                <Alert onClose={handleClose} severity="success">
                    Created shortened link at {path}
                </Alert>
            </Snackbar>
        </SnackbarProvider>
    )
}