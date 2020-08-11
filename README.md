# Uptime server API
[![Build Status](https://travis-ci.com/Marvin9/uptime-server-microservice.svg?token=VLAzbJP7VasfzqzWUHz9&branch=master)](https://travis-ci.com/Marvin9/uptime-server-microservice)

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
            <pre>
            {
                "email": string,
                "password": string,
            }
            </pre>
        </td>
        <td>
            <pre>
            {
                "email": bool,
                "password": string,
            }
            </pre>
        </td>
    </tr>
    <tr>
        <td>/auth/login</td>
        <td>POST</td>
        <td>
            <pre>
            {
                "email": string,
                "password": string
            }   
            </pre>
        </td>
        <td>
        <pre>
        {
            "error": bool,
            "message": string
        }
        </pre>
        </td>
        <td>
            Set jsonwebtoken in http cookie with expiration time of 30 minutes.
        </td>
    </tr>
    <tr>
        <td>/api/instance</td>
        <td>POST</td>
        <td>
            <pre>
            {
                "url": string,
                "duration": int (time in nanoseconds)
            }
            </pree>
        </td>
        <td>
            <pre>
            {
                "error": bool,
                "message": string
            }
            </pre>
        </td>
        <td>Set monitor for url which will check in duration period and update database accordingly. Reflects in /api/report</td>
    </tr>
    <tr>
        <td>/api/instances</td>
        <td>GET</td>
        <td>
            -
        </td>
        <td>
            <pre>
            {
                "error": bool,
                "data": {
                    "url": string,
                    "duration": int,
                    "unique_id": string,
                } || "message": "Unauthorized"
            }
            </pre>
        </td>
        <td>-</td>
    </tr>
    <tr>
        <td>/api/report</td>
        <td>GET</td>
        <td>
            -
        </td>
        <td>
            <pre>
            {
                "error": bool,
                "data": {
                    "url": string,
                    "status": int,
                    "reported_at": time,
                    "instance_id": string,
                } || "message": "Unauthorized"
            }
            </pre>
        </td>
        <td>-</td>
    </tr>
    <tr>
        <td>/api/instance</td>
        <td>DELETE</td>
        <td>
            <pre>
            {
                "instance_id": string,
            }
            </pre>
        </td>
        <td>
            <pre>
            {
                "error": bool,
                "message": string
            }
            </pre>
        </td>
        <td>-</td>
    </tr>
</table>