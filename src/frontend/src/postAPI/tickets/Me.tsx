import axiosBackend from "..";
import { UserInfo } from "types/Privilege";

interface resp {
    status: number
    data: UserInfo
}

const ProfileInfo = async function(): Promise<resp> {
    let response = await axiosBackend.get(`/me`);
    return response;
}
export default ProfileInfo;
