import axios from "axios";
import { FetchResponse, IContainer } from "../../types/types.ts";

const SearchContainers = async (query: string): Promise<IContainer[]> => {
    const response = await axios.get<FetchResponse>(`http://localhost:8080/api/v1/search?query=${query}`);
    return response.data.containers;  
};

export default SearchContainers;
