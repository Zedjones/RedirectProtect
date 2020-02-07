import React, { useState } from 'react';
import RedirectionFormatText from './RedirectionFormatText';
import RedirectionFormTime from './RedirectionFormTime';


export default function RedirectionForm(props) {
    const setURL = props.setURL;
    const setPassphrase = props.setPassphrase;
    const loading = props.loading;
    const durationEnabled = props.durationEnabled;
    const handleDurationEnableChange = props.handleDurationEnableChange;
    const selectedDate = props.selectedDate;
    const handleDateChange = props.handleDateChange;

    <form className={classes.form} noValidate onSubmit={(e) => createShortened(e)}>
        <RedirectionFormatText setURL={setURL} setPassphrase={setPassphrase} />
        <RedirectionFormTime
            loading={loading}
            durationEnabled={durationEnabled}
            handleDurationEnableChange={handleDurationEnableChange}
            selectedDate={selectedDate}
            handleDateChange={handleDateChange}
        />
    </form>
}