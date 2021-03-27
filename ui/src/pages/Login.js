import  { useState } from 'react'
import FormErrors from '../components/FormErrors';

export default function Login() {
  return(
    <div className="container">
      <div className="row">
        <div className="column column-50 column-offset-25">
          <h1>Login</h1>
          <h4>Enter your username and password to access your account</h4>
          <LoginForm />
        </div>
      </div>
    </div>
  );
}

function LoginForm() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [formErrors, setFormErrors] = useState([]);

  function handleFormInput(e) {
    switch(e.target.id) {
      case 'email':
        setEmail(e.target.value);
        break;
      case 'password':
        setPassword(e.target.value);
        break;
    }
  }

  function submitForm() {
    fetch('/login', {
      method: 'POST',
      body: JSON.stringify({
        "Email": email,
        "Password": password,
      })
    }).then(response => response.json())
      .then(payload => {
        if(payload.Error) {
          setFormErrors([payload.Error]);
        } else {
          handleLoginSuccess(payload);
        }
      });
  }

  function handleLoginSuccess() {

  }

  return(
    <fieldset>
      <FormErrors errors={formErrors} />

      <label htmlFor="email">E-mail Address</label>
      <input type="email" id="email" value={email} onChange={handleFormInput} />

      <label htmlFor="password">Password</label>
      <input type="password" id="password" value={password} onChange={handleFormInput} />

      <button className="input-primary" onClick={submitForm}>Login</button>
    </fieldset>
  )
}
