import React from 'react';
import TextField from '@material-ui/core/TextField';

export default function RedirectionFormText(props) {
    const [setURL, setPassphrase] = [props.setURL, props.setPassphrase];

    return (
        <div>
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
        </div>
    )
}