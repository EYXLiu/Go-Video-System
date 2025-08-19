//Chunking is done client side, as not to fail and restart from 0% AND to not lose any progress + less load on server

//requires to be in a ReactJS Application (for process.env)
const CHUNK_SIZE = 5 * 1024 * 1024;

const initUpload = async (file) => {
    const response = await fetch(`${process.env.URL}/upload/init`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            fileName: file.name,
            totalChunks: Math.ceil(file.size / CHUNK_SIZE)
        })
    });
    const data = await response.json();
    return data.uploadID;
}

const uploadChunks = async(file, uploadID) => {
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

    for (let i = 0; i < totalChunks; i++) {
        const start = i * CHUNK_SIZE;
        const end = Math.min(file.size, start + CHUNK_SIZE);
        const chunk = file.slice(start, end);

        const formData = new FormData();
        formData.append("uploadID", uploadID);
        formData.append("chunkNum", i);
        formData.append("file", chunk);

        await fetch(`${process.env.URL}/upload/chunk`, {
            method: "POST", 
            body: formData
        });
        //log if required
    }
}

const uploadComplete = async (uploadID) => {
    const formData = new FormData;
    formData.append("uploadID", uploadID);

    const response = await fetch(`${process.env.URL}/upload/complete`, {
        method: "POST", 
        body: formData
    });

    const data = await response.json();
    //log if required
}

const uploadVideo = async (file) => {
    const uploadID = await initUpload(file);
    await uploadChunks(file, uploadID);
    await uploadComplete(uploadID);
}

const getVideo = async(id, res) => {
    const response = await fetch(`${process.env.URL}/video/${id}`, {
        method: "GET"
    }) 
    const data = await response.json();
    const thumbnail = data.video.Resolutions["thumbnail"];
    const url = data.video.Resolutions[res];
}