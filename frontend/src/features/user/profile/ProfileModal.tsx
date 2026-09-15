import {
    useEffect,
    useRef,
    useState,
    type ChangeEvent,
    type FormEvent,
} from 'react'

import {
    Camera,
    LoaderCircle,
    X,
} from 'lucide-react'

import {
    updateProfile,
    uploadAvatar,
    type UserProfile,
} from './profileApi'

import './ProfileModal.css'

type ProfileModalProps = {
    profile: UserProfile
    onClose: () => void
    onUpdated: (profile: UserProfile) => void
}

function ProfileModal({
                          profile,
                          onClose,
                          onUpdated,
                      }: ProfileModalProps) {
    const fileInputRef =
        useRef<HTMLInputElement>(null)

    const [name, setName] =
        useState(profile.name)

    const [username, setUsername] =
        useState(profile.username)

    const [bio, setBio] =
        useState(profile.bio ?? '')

    const [avatarFile, setAvatarFile] =
        useState<File | null>(null)

    const [avatarPreview, setAvatarPreview] =
        useState<string | null>(
            profile.avatar_path || null,
        )

    const [isSaving, setIsSaving] =
        useState(false)

    const [error, setError] =
        useState('')

    useEffect(() => {
        if (!avatarFile) {
            return
        }

        const previewUrl =
            URL.createObjectURL(avatarFile)

        setAvatarPreview(previewUrl)

        return () => {
            URL.revokeObjectURL(previewUrl)
        }
    }, [avatarFile])

    function handleAvatarChange(
        event: ChangeEvent<HTMLInputElement>,
    ) {
        const file =
            event.target.files?.[0]

        if (!file) {
            return
        }

        setError('')

        if (file.size > 5 * 1024 * 1024) {
            setError(
                'Avatar must not exceed 5 MB.',
            )
            return
        }

        setAvatarFile(file)
    }

    async function handleSubmit(
        event: FormEvent<HTMLFormElement>,
    ) {
        event.preventDefault()

        if (isSaving) {
            return
        }

        setIsSaving(true)
        setError('')

        try {
            const profileResponse =
                await updateProfile({
                    name,
                    username,
                    bio,
                })

            if (!profileResponse.success) {
                setError(
                    profileResponse.error.message,
                )
                return
            }

            let updatedProfile =
                profileResponse.data

            if (avatarFile) {
                const avatarResponse =
                    await uploadAvatar(avatarFile)

                if (!avatarResponse.success) {
                    setError(
                        avatarResponse.error.message,
                    )
                    return
                }

                updatedProfile =
                    avatarResponse.data
            }

            onUpdated(updatedProfile)
            onClose()
        } catch {
            setError(
                'Could not update profile.',
            )
        } finally {
            setIsSaving(false)
        }
    }

    function getInitials() {
        return name
            .split(' ')
            .map((part) =>
                part.charAt(0),
            )
            .join('')
            .slice(0, 2)
            .toUpperCase()
    }

    return (
        <div className="profile-modal-backdrop">
            <div className="profile-modal">
                <div className="profile-modal-header">
                    <div>
                        <h2>Edit profile</h2>
                        <p>
                            Update your personal
                            information.
                        </p>
                    </div>

                    <button
                        type="button"
                        className="profile-close-button"
                        onClick={onClose}
                    >
                        <X size={20} />
                    </button>
                </div>

                <form
                    className="profile-form"
                    onSubmit={handleSubmit}
                >
                    <div className="profile-avatar-section">
                        <div className="profile-avatar">
                            {avatarPreview ? (
                                <img
                                    src={
                                        avatarPreview
                                    }
                                    alt={name}
                                />
                            ) : (
                                <span>
                                    {getInitials()}
                                </span>
                            )}
                        </div>

                        <button
                            type="button"
                            className="profile-avatar-button"
                            onClick={() =>
                                fileInputRef.current?.click()
                            }
                        >
                            <Camera size={17} />
                            Change avatar
                        </button>

                        <input
                            ref={fileInputRef}
                            type="file"
                            accept="image/jpeg,image/png,image/webp"
                            hidden
                            onChange={
                                handleAvatarChange
                            }
                        />
                    </div>

                    <div className="profile-field">
                        <label htmlFor="profile-name">
                            Name
                        </label>

                        <input
                            id="profile-name"
                            type="text"
                            value={name}
                            maxLength={100}
                            onChange={(event) =>
                                setName(
                                    event.target.value,
                                )
                            }
                        />
                    </div>

                    <div className="profile-field">
                        <label htmlFor="profile-username">
                            Username
                        </label>

                        <div className="profile-username-input">
                            <span>@</span>

                            <input
                                id="profile-username"
                                type="text"
                                value={username}
                                maxLength={50}
                                onChange={(event) =>
                                    setUsername(
                                        event.target.value,
                                    )
                                }
                            />
                        </div>
                    </div>

                    <div className="profile-field">
                        <div className="profile-field-header">
                            <label htmlFor="profile-bio">
                                Bio
                            </label>

                            <span>
                                {bio.length}/160
                            </span>
                        </div>

                        <textarea
                            id="profile-bio"
                            value={bio}
                            maxLength={160}
                            rows={4}
                            placeholder="Tell people a little about yourself..."
                            onChange={(event) =>
                                setBio(
                                    event.target.value,
                                )
                            }
                        />
                    </div>

                    {error && (
                        <div className="profile-error">
                            {error}
                        </div>
                    )}

                    <div className="profile-actions">
                        <button
                            type="button"
                            className="profile-cancel-button"
                            onClick={onClose}
                            disabled={isSaving}
                        >
                            Cancel
                        </button>

                        <button
                            type="submit"
                            className="profile-save-button"
                            disabled={
                                isSaving ||
                                !name.trim() ||
                                !username.trim()
                            }
                        >
                            {isSaving ? (
                                <>
                                    <LoaderCircle
                                        size={17}
                                        className="profile-spinner"
                                    />
                                    Saving...
                                </>
                            ) : (
                                'Save changes'
                            )}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default ProfileModal