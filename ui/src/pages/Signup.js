import { useState } from 'react';
import { Redirect } from "react-router-dom";

import FormErrors from '../components/FormErrors';

export default function Signup() {
  const [redirectRoot, setRedirectRoot] = useState(false);

  function handleSignupSuccess() {
    setRedirectRoot(true);
  }

  if(redirectRoot) {
    return(<Redirect to='/' />);
  }

  return(
    <div className="container">
      <div className="row">
        <div className="column column-50 column-offset-25">
          <h1>Signup</h1>
          <h4>Enter your e-mail address, registration number and a password to get automatic reminders when your vehicle is due an MOT.</h4>
          <SignupForm onSuccess={handleSignupSuccess} />
        </div>
      </div>
    </div>
  );
}

function SignupForm({ onSuccess }) {
  const [email, setEmail] = useState("");
  const [registrationNumber, setRegistrationNumber] = useState("")
  const [password, setPassword] = useState("");
  const [passwordConfirm, setPasswordConfirm] = useState("");
  const [termsAndConditions, setTermsAndConditions] = useState(false);
  const [formErrors, setFormErrors] = useState([]);

  function handleFormInput(e) {
    switch(e.target.id) {
      case 'email':
        setEmail(e.target.value);
        break;
      case 'registrationNumber':
        setRegistrationNumber(e.target.value);
        break;
      case 'password':
        setPassword(e.target.value);
        break;
      case 'passwordConfirmation':
        setPasswordConfirm(e.target.value);
        break;
      case 'termsAndConditions':
        setTermsAndConditions(e.target.checked)
        break;
    }
  }

  function submitForm() {
    fetch('/signup', {
      method: 'POST',
      body: JSON.stringify({
        "Email": email,
        "RegistrationNumber" : registrationNumber,
        "Password": password,
        "TermsAndConditions": termsAndConditions
      })
    }).then(response => response.json())
      .then(payload => {
        if(payload.Error) {
          setFormErrors(payload.Error);
        } else {
          loginUser();
        }
      });
  }

  function loginUser() {
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
          onSuccess();
        }
      });
  }

  return(
    <fieldset>
      <FormErrors errors={formErrors} />

      <label htmlFor="email">E-mail Address</label>
      <input type="email" id="email" value={email} onChange={handleFormInput} />

      <label htmlFor="registrationNumber">Vehicle Registration Number</label>
      <input type="text" id="registrationNumber" value={registrationNumber} onChange={handleFormInput} />

      <label htmlFor="password">Password</label>
      <input type="password" id="password" value={password} onChange={handleFormInput} />

      <label htmlFor="password">Confirm Password</label>
      <input type="password" id="passwordConfirmation" value={passwordConfirm} onChange={handleFormInput} />

      <div>
        <input type="checkbox" id="termsAndConditions" checked={termsAndConditions} onChange={handleFormInput} />
        <label className="label-inline" htmlFor="termsAndConditions">I agree to the terms and conditions and privacy policy</label>
      </div>

      <button className="input-primary" onClick={submitForm}>Complete Signup</button>
    </fieldset>
  )
}
