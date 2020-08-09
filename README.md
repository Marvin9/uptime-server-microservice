# Uptime server API

> Monitor your server uptime easily.

## API

<table>
    <tr>
        <th>Endpoint</th>
        <th>Method</th>
        <th>Request Body</th>
        <th>Response</th>
        <th>Additional</th>
    </tr>
    <tr>
        <td>/auth/register</td>
        <td>POST</td>
        <td>
            <code>
                "email": string,
                "password": string,
            </code>
        </td>
        <td>
            <code>
                "email": bool,
                "password": string,
            </code>
        </td>
    </tr>
    <tr>
        <td>/auth/login</td>
        <td>POST</td>
        <td><code>
        "email": string,
        "password": string
        </code></td>
        <td>
        <code>
        "error": bool,
        "message": string
        </code>
        </td>
        <td>
            Set jsonwebtoken in http cookie with expiration time of 30 minutes.
        </td>
    </tr>
    <tr>
        <td>/api/instance</td>
        <td>POST</td>
        <td>
            <code>
                "url": string,
                "duration": int (time in nanoseconds)
            </code>
        </td>
        <td>
            <code>
                "error": bool,
                "message": string
            </code>
        </td>
        <td>Set monitor for url which will check in duration period and update database accordingly. Reflects in /api/report</td>
    </tr>
    <tr>
        <td>/api/report</td>
        <td>GET</td>
        <td>
            -
        </td>
        <td>
            <code>
                "error": bool,
                "data": {
                    "url": string,
                    "status": int,
                    "reported_at": time,
                    "instance_id": string,
                } OR
                "message": "Unauthorized"
            </code>
        </td>
        <td>-</td>
    </tr>
    <tr>
        <td>/api/instance</td>
        <td>DELETE</td>
        <td>
            <code>
                "instance_id": string,
            </code>
        </td>
        <td>
            <code>
                "error": bool,
                "message": string
            </code>
        </td>
        <td>-</td>
    </tr>
</table>