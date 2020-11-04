<?php

namespace App\Models;

/**
 * Interface PasswordHasher
 *
 * @package App\Models
 */
interface PasswordHasher
{
    /**
     * @param string $password
     *
     * @return string
     */
    public function hash(string $password): string;

    /**
     * @param string $password
     * @param string $hashedPassword
     *
     * @return bool
     */
    public function check(string $password, string $hashedPassword): bool;
}
