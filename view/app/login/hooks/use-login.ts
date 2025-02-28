import React from 'react'

function useLogin() {
    const [email, setEmail] = React.useState('')
    const [password, setPassword] = React.useState('')

    const handleEmailChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmail(event.target.value)
    }

    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPassword(event.target.value)
    }

    const handleLogin = () => {
        console.log(email, password)
    }

    return {
        email,
        password,
        handleEmailChange,
        handlePasswordChange,
        handleLogin
    }
}

export default useLogin