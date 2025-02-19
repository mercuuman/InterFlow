import '../style/registration.css';
import {registrationUser} from './components/api'

const registrationForm = document.forms.registration_user;
const registrationFormNameInput = registrationForm.elements.user_name;
const registrationFormMailInput = registrationForm.elements.user_mail;
const registrationFormPasswordInput = registrationForm.elements.user_password;

const registrationMessage = document.querySelector('.registration__message');
const registrationMessageLink = document.querySelector('.message__link');

function registrationHandler(evt) {
  evt.preventDefault();
  const registrationFormMailInput = registrationForm.elements.user_mail;
  registrationUser({name: registrationFormNameInput.value, mail: registrationFormMailInput.value, password: registrationFormPasswordInput.value})
  .then((res) => {
    registrationForm.closest('.registration__content').classList.add('hidden');
    registrationMessage.classList.add('registration__message_visible');
  }).catch(err => {
    console.log(err);
  })
}

registrationForm.addEventListener('submit', registrationHandler)
