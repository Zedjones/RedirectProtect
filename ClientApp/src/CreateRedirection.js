import React, { useState } from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import Fade from '@material-ui/core/Fade';
import CircularProgress from '@material-ui/core/CircularProgress';
import TextField from '@material-ui/core/TextField';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import Checkbox from '@material-ui/core/Checkbox';
import { makeStyles } from '@material-ui/core/styles';
import { KeyboardTimePicker } from "@material-ui/pickers";
import Container from '@material-ui/core/Container';
import lifecycle from 'react-pure-lifecycle';
import { MuiPickersUtilsProvider } from "@material-ui/pickers";
import LuxonUtils from "@date-io/luxon"
import { DateTime } from "luxon"
import AccessTimeIcon from "@material-ui/icons/AccessTime"
import { Grid, List, ListItem, ListItemText } from '@material-ui/core';

const useStyles = makeStyles(theme => ({
  '@global': {
    body: {
      backgroundColor: theme.palette.grey,
    },
  },
  paper: {
    marginTop: theme.spacing(8),
    display: "flex",
    alignItems: "center",
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main,
  },
  form: {
    width: '100%', // Fix IE 11 issue.
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
  list: {
    marginTop: theme.spacing(7.5)
  }
}));

const methods = {
  componentDidMount() {
    document.title = "Redirect Protect"
  }
}

function SignIn() {
  const classes = useStyles();
  const [selectedDate, handleDateChange] = useState(new DateTime.fromObject({ hours: 0, minutes: 0 }));
  const [durationEnabled, handleDurationEnableChange] = useState(false);
  const [URL, setURL] = useState('');
  const [loading, setLoading] = React.useState(false);
  const [passphrase, setPassphrase] = useState('');

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
        if (response.status === 200) {
        }
        else {
          response.json().then(json => console.log(json));
        }
        setLoading(false);
      })
      .catch(((_) => {
        setLoading(false);
      }))
  };

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <Grid container className={classes.paper} direction="column">
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Shorten & Encrypt URL
          </Typography>
        <form className={classes.form} noValidate onSubmit={(e) => createShortened(e)}>
          <TextField
            variant="outlined"
            margin="normal"
            required
            disabled={loading}
            fullWidth
            id="url"
            label="URL"
            name="url"
            onChange={(ev) => setURL(ev.target.value)}
            autoComplete="url"
            type="url"
            autoFocus
          />
          <TextField
            variant="outlined"
            margin="normal"
            required
            disabled={loading}
            fullWidth
            name="password"
            label="Passphrase"
            type="password"
            id="password"
            onChange={(ev) => setPassphrase(ev.target.value)}
            autoComplete="current-password"
          />
          <Grid container direction='row' alignItems='center'>
            <Grid item xs={2}>
              <Checkbox
                checked={durationEnabled}
                disabled={loading}
                onChange={(event) => handleDurationEnableChange(event.target.checked)}
                inputProps={{
                  'aria-label': 'primary checkbox',
                }}
              />
            </Grid>
            <Grid item xs={10}>
              <MuiPickersUtilsProvider utils={LuxonUtils}>
                <KeyboardTimePicker
                  clearable
                  ampm={false}
                  disabled={!durationEnabled || loading}
                  margin="normal"
                  autoOk={true}
                  views={["hours", "minutes"]}
                  inputVariant="outlined"
                  label="Shortened URL Lifespan"
                  openTo="minutes"
                  format="HH:mm"
                  keyboardIcon={React.createElement(AccessTimeIcon, null)}
                  placeholder="hh:mm"
                  value={selectedDate}
                  style={{ width: "100%" }}
                  onChange={(val) => { handleDateChange(val); console.log(selectedDate) }}
                />
              </MuiPickersUtilsProvider>
            </Grid>
          </Grid>
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
    </Container >
  );
}

export default lifecycle(methods)(SignIn);