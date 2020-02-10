import React, { useState, Fragment } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Fade from '@material-ui/core/Fade';
import { DateTime } from "luxon"
import CircularProgress from '@material-ui/core/CircularProgress';
import { Grid, IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';

import RedirectionFormText from './RedirectionFormText';
import RedirectionFormTime from './RedirectionFormTime';
import { useSnackbar } from 'notistack';

const useStyles = makeStyles(theme => ({
    form: {
        width: '100%', // Fix IE 11 issue.
        marginTop: theme.spacing(1),
    },
    submit: {
        margin: theme.spacing(3, 0, 2),
    },
    paper: {
        display: "flex",
        alignItems: "center",
    },
}));

export default function RedirectionForm(props) {
    const classes = useStyles();
    const [URL, setURL] = useState('');
    const [passphrase, setPassphrase] = useState('');
    const [durationEnabled, handleDurationEnableChange] = useState(false);
    const [loading, setLoading] = useState(false);
    const [selectedDate, handleDateChange] = useState(new DateTime.fromObject({ hours: 0, minutes: 0 }));
    const { enqueueSnackbar, closeSnackbar } = useSnackbar();
    const onClickDismiss = key => () => {
        closeSnackbar(key);
    }

    const action = key => (
        <React.Fragment>
            <IconButton
                size="small"
                aria-label="close"
                color="inherit"
                onClick={onClickDismiss(key)}
            >
                <CloseIcon fontSize="small" />
            </IconButton>
        </React.Fragment>
    )

    function createShortened(ev) {
        ev.preventDefault();
        let ttl = null;
        let missing = [];
        URL === "" ? missing.push("URL") : null;
        passphrase === "" ? missing.push("passphrase") : null;
        if (missing.length !== 0) {
            enqueueSnackbar(`Please include a ${missing.join(' and ')}.`, {
                variant: 'warning',
                action
            });
            return;
        }
        if (durationEnabled && selectedDate.c != null) {
            ttl = selectedDate.toISO()
        }
        let redirect = {
            "URL": URL,
            "Password": passphrase,
            "TTL": ttl
        }

        setLoading(true);
        fetch('api/redirect', {
            method: 'POST',
            body: JSON.stringify(redirect),
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            }
        })
            .then(response => {
                setLoading(false);
                if (response.status === 200) {
                    response.text().then(text => {
                        enqueueSnackbar(`Created shortened link at ${text}`, {
                            variant: 'success',
                            action
                        });
                    });
                }
                else {
                    response.json().then(json => console.log(json));
                }
            })
            .catch(((_) => {
                setLoading(false);
            }))
    };

    return (
        <Grid container className={classes.paper} direction="column">
            <form className={classes.form} noValidate onSubmit={(e) => createShortened(e)}>
                <RedirectionFormText setURL={setURL} setPassphrase={setPassphrase} loading={loading} />
                <RedirectionFormTime
                    loading={loading}
                    durationEnabled={durationEnabled}
                    handleDurationEnableChange={handleDurationEnableChange}
                    selectedDate={selectedDate}
                    handleDateChange={handleDateChange}
                />
                <Button
                    type="submit"
                    fullWidth
                    variant="contained"
                    color="primary"
                    className={classes.submit}
                    disabled={loading}
                >
                    Shorten
                </Button>
            </form>
            <Fade
                in={loading}
                style={{
                    transitionDelay: loading ? '800ms' : '0ms',
                }}
                unmountOnExit
            >
                <CircularProgress />
            </Fade>
        </Grid>
    );
}