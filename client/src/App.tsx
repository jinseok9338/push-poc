import React, { useEffect, useState } from 'react';


import AWS from 'aws-sdk';
import { Alert, Snackbar } from '@mui/material';

const SNS_TOPIC_ARN = 'arn:aws:sns:ap-northeast-2:117371106642:notification';
const REGION = 'ap-northeast-2';

console.log(import.meta.env);
AWS.config.update({
  region: REGION,
  credentials: {
    accessKeyId: import.meta.env.VITE_AWS_ACCESS_KEY_ID!,
    secretAccessKey: import.meta.env.VITE_AWS_SECRET_ACCESS_KEY!,
  }
});

const sns = new AWS.SNS();

function App() {
  const [message, setMessage] = useState('');
  const [open, setOpen] = useState(false);
  const [severity, setSeverity] = useState<'success' | 'info' | 'warning' | 'error'>('info');

  useEffect(() => {
    sns.subscribe(
      {
        Protocol: 'http',
        TopicArn: SNS_TOPIC_ARN,
        Endpoint: "http://8dcd-221-143-168-170.ngrok.io/" ,
      },
      (err, data) => {
        if (err) {
          console.error(err);
        } else {
          console.log(data);
        }
      }
    );

    const receiveNotification = (event:any) => {
      if (event.data) {
        const notification = JSON.parse(event.data);

        if (notification.Message) {
          setMessage(notification.Message);
          setSeverity(notification.severity || 'info');
          setOpen(true);
        }
      }
    };

    window.addEventListener('message', receiveNotification);

    return () => {
      window.removeEventListener('message', receiveNotification);
    };
  }, []);

  return (
    <div>
      <Snackbar open={open} autoHideDuration={6000} onClose={() => setOpen(false)}>
        <Alert onClose={() => setOpen(false)} severity={severity}>
          {message}
        </Alert>
      </Snackbar>
      <h1>Subscribe to SNS Topic</h1>
    </div>
  );
}

export default App;




