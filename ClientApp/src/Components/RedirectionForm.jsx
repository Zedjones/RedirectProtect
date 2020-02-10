import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Fade from '@material-ui/core/Fade';
import { DateTime } from "luxon"
import CircularProgress from '@material-ui/core/CircularProgress';
import { Grid } from '@material-ui/core';

import RedirectionFormText from './RedirectionFormText';
import RedirectionFormTime from './RedirectionFormTime';
import CreateToast from '../CreateToast';

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
    const [toastOpen, setToastOpen] = useState(false);
    const [lastPath, setLastPath] = useState('');

    function createShortened(ev) {
        ev.preventDefault();
        let ttl = null;
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
                        //TODO: update this for our new enqueueSnackbar approach
                        setLastPath(text);
                        setToastOpen(true);
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
                <CreateToast open={toastOpen} setOpen={setToastOpen} path={lastPath}></CreateToast>
                <Button
                    type="submit"
                    fullWidth
                    variant="contained"
                    color="primary"
                    className={classes.submit}
                    disabled={loading}
                >
                    Encrypt
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