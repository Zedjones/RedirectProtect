import { SnackbarProvider } from "notistack";
import React from 'react';
import Snackbar from '@material-ui/core/Snackbar';
import MuiAlert from '@material-ui/lab/Alert';
import CloseIcon from '@material-ui/icons/Close';
import { IconButton } from "@material-ui/core";

function Alert(props) {
    return <MuiAlert elevation={0} variant="filled" {...props} />
}

export default function CreateToast(props) {
    const [open, setOpen] = [props.open, props.setOpen];
    const path = props.path
    const notistackRef = React.createRef();
    const onClickDismiss = key => () => {
        notistackRef.current.closeSnackbar(key);
    }

    const handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }

        setOpen(false);
    };

    return (
        <SnackbarProvider
            ref={notistackRef}
            action={(key) => (
                <React.Fragment>
                    <IconButton size="small" aria-label="close" color="inherit" onClick={onClickDismiss(key)}>
                        <CloseIcon fontSize="small" />
                    </IconButton>
                </React.Fragment>
            )}
            anchorOrigin={{ horizontal: "center" }}>
            <Snackbar open={open} autoHideDuration={4000} onClose={handleClose}>
                <Alert onClose={handleClose} severity="success">
                    Created shortened link at {path}
                </Alert>
            </Snackbar>
        </SnackbarProvider>
    )
}