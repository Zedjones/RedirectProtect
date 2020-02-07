import React from 'react';
import LuxonUtils from "@date-io/luxon";
import { Grid } from '@material-ui/core';
import Checkbox from '@material-ui/core/Checkbox';
import { MuiPickersUtilsProvider } from "@material-ui/pickers";
import AccessTimeIcon from "@material-ui/icons/AccessTime";
import { KeyboardTimePicker } from "@material-ui/pickers";

export default function RedirectionFormTime(props) {
    const loading = props.loading;
    const durationEnabled = props.durationEnabled;
    const handleDurationEnableChange = props.handleDurationEnableChange;
    const selectedDate = props.selectedDate;
    const handleDateChange = props.handleDateChange;

    return (
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
    );
}