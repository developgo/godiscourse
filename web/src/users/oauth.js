import LoadingPage from '../loading/loading_page.js';
import URLUtils from '../components/url.js';
import React, { Component } from 'react';
import API from '../api/index.js';

class Oauth extends Component {
  constructor(props) {
    super(props);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('loading', 'layout');
    this.api = new API();
  }

  componentDidMount() {
    const history = this.props.history;
    const code = new URLUtils().getUrlParameter('code');
    this.api.user.signIn(code).then((resp) => {
      history.push('/');
    });
  }

  render() {
    return (
      <LoadingPage />
    );
  }
}

export default Oauth;
