import axiosBackend from "..";
import { AuthRequest } from "types/Account";

interface resp {
  status: number
}

export const Login = async function (data: AuthRequest): Promise<resp> {
  const response = await axiosBackend
    .post(`/authorize`, data, { withCredentials: true })
    .then((data) => data)
    .catch((error) => {
      return { status: error.response?.status, data: error.response.data };
    });

  if (response.status === 200) {
    localStorage.setItem("authToken", response.data.access_token);
    console.log(response.data.role)
    localStorage.setItem("role", response.data.role);
  //localStorage.setItem("role", "user");

  }

  return {
    status: response?.status,
  };
};
