import React from 'react';
import Avatar from '@material-ui/core/Avatar';
import CssBaseline from '@material-ui/core/CssBaseline';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import lifecycle from 'react-pure-lifecycle';
import { Grid } from '@material-ui/core';
import RedirectionForm from './Components/RedirectionForm';
import { SnackbarProvider } from "notistack";

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
        <SnackbarProvider>
          <RedirectionForm />
        </SnackbarProvider>
      </Grid>
    </Container >
  );
}

export default lifecycle(methods)(SignIn);