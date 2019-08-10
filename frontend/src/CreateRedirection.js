import React, { useState } from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import { KeyboardTimePicker } from "@material-ui/pickers";
import Container from '@material-ui/core/Container';
import lifecycle from 'react-pure-lifecycle';
import { MuiPickersUtilsProvider } from "@material-ui/pickers";
import LuxonUtils from "@date-io/luxon"
import { DateTime } from "luxon"
import AccessTimeIcon from "@material-ui/icons/AccessTime"

const useStyles = makeStyles(theme => ({
  '@global': {
    body: {
      backgroundColor: theme.palette.grey,
    },
  },
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
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
}));

const methods = {
  componentDidMount() {
    document.title = "Redirect Protect"
  }
}

function SignIn() {
  const classes = useStyles();
  const [selectedDate, handleDateChange] = useState(new DateTime.fromObject({ hours: 0, minutes: 0 }));

  return (
    <MuiPickersUtilsProvider utils={LuxonUtils}>
      <Container component="main" maxWidth="xs">
        <CssBaseline />
        <div className={classes.paper}>
          <Avatar className={classes.avatar}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Shorten & Encrypt URL
        </Typography>
          <form className={classes.form} noValidate>
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              id="username"
              label="URL"
              name="username"
              autoComplete="username"
              type="url"
              autoFocus
            />
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              name="password"
              label="Passphrase"
              type="password"
              id="password"
              autoComplete="current-password"
            />
            <KeyboardTimePicker
              clearable
              ampm={false}
              autoOk={true}
              views={["hours", "minutes"]}
              label="Shortened URL Lifespan"
              openTo="minutes"
              format="HH:mm"
              keyboardIcon={React.createElement(AccessTimeIcon, null)}
              placeholder="00:10"
              value={selectedDate}
              style={{width: "100%"}}
              onChange={(val) => { handleDateChange(val); console.log(selectedDate) }}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
            >
              Encrypt
          </Button>
          </form>
        </div>
      </Container>
    </MuiPickersUtilsProvider>
  );
}

export default lifecycle(methods)(SignIn);